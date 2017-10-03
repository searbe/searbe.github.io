Notes on how to use OpenSSL and your own CA to sign certificates. This can be very useful
when you want to issue yourself multiple certificates, but don't want to pin against each
of those certificates individually. Instead, trust your own CA and use it to sign your 
certificates.

Once it's configured you'll be able do this to create certificates:

    $ openssl req -new -key ~/my/key -out request.csr
    $ openssl ca -in request.csr -out cert.pem
    $ ls
    cert.pem
    request.csr

# OpenSSL

Man page:

    $ man openssl

> TL;DR: "OpenSSL does all sorts of cryptography stuff, not just SSL."

At the bottom of the man page you'll find references to other man pages. These are all man
pages for OpenSSL sub-commands - to run the command described in `man ca` you need to run
`openssl ca`, and so-on.

# Configure your Certificate Authority ("CA")

Man page:

    $ man ca

> TL;DR: "Used to sign certificate requests, generate CRLs [revoked certificate lists] and maintain a database of issued
certificates." You will also find documentation on CA configuration.

On Ubuntu you'll find your CA configuration file already exists at `/etc/ssl/openssl.cnf`.
Pay attention to this configuration, especially:

    dir         = /etc/ssl
    certificate = $dir/cacert.pem
    private_key = $dir/private/cakey.pem

You'll want to make sure the following configuration references things which exist:

    new_certs_dir = $dir/newcerts
    database      = $dir/index.txt
    serial        = $dir/serial

    $ sudo mkdir /etc/ssl/newcerts
    $ sudo touch /etc/ssl/index.txt
    $ echo "01" | sudo tee /etc/ssl/serial

# Policies

Note this in the configuration:

    policy = policy_match

When you create your CA's root certificate (next step), all certificate requests must adhere
to your policy. You can configure multiple policies - you might have some defaults:

    [ policy_match ]
    countryName             = match
    stateOrProvinceName     = match
    organizationName        = match
    organizationalUnitName  = optional
    commonName              = supplied
    emailAddress            = optional

    [ policy_anything ]
    countryName             = optional
    stateOrProvinceName     = optional
    localityName            = optional
    organizationName        = optional
    organizationalUnitName  = optional
    commonName              = supplied
    emailAddress            = optional

This default `policy_match` policy will require that any certificate request has
the same countryName, stateOrProvinceName and organizationName as the CA's root certificate.

To make your habits good by default I suggest that you create a new policy - e.g.
`policy_match_organisation` - and configure it how you like, setting it as the default policy.
I like to change `policy_anything` so commonName is optional too.

# Create your root certificate

Man page:

    $ man req

> TL;DR: "Creates and processes certificate requests."

You'll find an example: "Generate a self signed root certificate."

    $ openssl req -x509 -newkey rsa:2048 -keyout key.pem -out req.pem

Look up the documentation for each argument:

    -x509
      outputs a self signed certificate
    -newkey
      this option creates a new certificate request and a new private key
    -keyout
      the filename to write the newly created private key to
    -out
      the output filename to write to 

Run it, but set the paths for the keys to what we saw in our configuration earlier. We need
to be root to use those paths.

    $ sudo openssl req -x509 -newkey rsa:2048 -keyout /etc/ssl/private/cakey.pem -out /etc/ssl/cacert.pem

Congratulations, you are now a certificate authority.

# Request a certificate

Now you want to start issuing certificates. Anybody who wants a certificate has to create a
certificate request.

## Create a key

This key identifies the entity requesting the certificate. You could use your ssh
key, even, but we'll generate a new one.

Man page:

    $ man genpkey

> TL;DR: "Generates private keys."

Following examples you'll see in the `man` page...

    $ openssl genpkey -algorithm RSA -out key.pem -pkeyopt rsa_keygen_bits:2048
    $ ls | grep key
    key.pem

## Create a certificate request

Use the key to generate a certificate request.

Man page:

    $ man req

> TL;DR: "Creates and processes certificate requests."

    $ openssl req -new -key key.pem -out request.csr
    [ ... answer the questions ... ]
    $ ls | grep csr
    request.csr

We have a certificate request.

# Process the request

Now we have a request for a certificate, we want to create the certificate.

Back to this man page:

    $ man ca

One example is "Sign a certificate request."

    $ openssl ca -in req.pem -out newcert.pem

Note this argument, elsewhere in the man page:

    -policy arg
        this option defines the CA "policy" to use

Sign the request with the 'policy_anything' policy we saw earlier. You'll need
root, as only root can read the CA key.

    $ sudo openssl ca -in request.csr -out newcert.pem -policy policy_anything
    $ ls | grep pem
    key.pem
    newcert.pem

You'll also find the new certificate in newcerts with serial number 01:

    $ ls /etc/ssl/newcerts
    01.pem

You'll see that the serial number has been set to the next one:

    $ cat /etc/ssl/serial
    02

You'll find an entry for your certificate in the index:

    $ cat /etc/ssl/index.txt
    V	180227092316Z		01	unknown	/C=UK/ST=Some-State/O=Internet Widgits Pty Ltd

You can grep out the certificate text from the `.pem` file:

    $ grep BEGIN newcert.pem -A 1000
    -----BEGIN CERTIFICATE-----
    MIIDgDCCAmigAwIBAgIBATANBgkqhkiG9w0BAQsFADBFMQswCQYDVQQGEwJBVTET
    ...
    txvPo+ilv5goqXo8FK6ghQzu8W6RtL1dmsL3sDy4pPFTJWre
    -----END CERTIFICATE-----

There you go. A certificate signed by your CA.