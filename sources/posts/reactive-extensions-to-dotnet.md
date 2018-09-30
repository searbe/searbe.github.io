This aims to get to the point with [Reactive Extensions to .Net](https://github.com/dotnet/reactive). If you want to drill down further, I suggest you read the fanstastic book at [IntroToRx](http://www.introtorx.com/).

If you get bored, at least read the summary and advice at the end.

## Observers

`IObservable<T>` is a stream of `T` which an `IObserver<T>` can be subscribed to, a-la `observable.Subscribe(observer)`.

Thus, an `IObserver<T>` observes instances of `T`.

An `IObservable<T>` will call `.OnNext()`, `.OnError()` or `.OnCompleted()` methods on each observer. When an observer receives `.OnComplete()` or `.OnError()`, it should expect no further errors or values.

There are overloads / extension methods for `IObservable<T>.Subscribe` which allow you to subscribe with a `(T t) => {}`, rather than an `IObserver<T>` implementation.

## Subjects

`Subject<T>` implements both interfaces. Anything pushed to it via `IObserver<T>.OnNext` is pushed to _it's own_ observers. Observers are registered through the subject's implementation of the `IObservable<T>` interface.

`Subject<T>` implements `ISubject<T>`, which itself implements `ISubject<TSource, TResult>` where both types are the same. That is; it is entirely possible for a subject to accept one type but publish another.

????? `Subject.Create()` accepts an `IObserver<TSource>` and an `IObservable<TResult>`. This is the easiest way to utilise a subject with different input and output types.

`ReplaySubject<T>` caches values; observers will receive events even if they subscribe _after_ events are published. A replay subject has a configurable `bufferSize` and a configurable time window.

`BehaviourSubject<T>` always contains a buffer of one `T` - so something which observes it will immediately get one value. When you construct a `BehaviourSubject<T>` you must provide the intitial `T`.

`AsyncSubject<T>` buffers one value but doesn't publish it until it's `.OnComplete()` method is called.

## Lifetime Management

`Subscribe()` will return an `IDisposable`. You should dispose it when you want to cancel the subscription. Consider utilizing C#'s `using (){}` statement.

There are overloads which accept cancellation tokens. Be careful not to end up with subscriptions which are never cancelled.

## Creating sequences

`Observable` has static methods which return `IObservable<T>` implementations:

    Observable.Return<string>("Value") // Returns one value to observer
    Observable.Empty<string>() // Calls OnComplete() instantly
    Observable.Never<string>() // Never completes
    Observable.Throw<string>(e) // Calls OnError(e) instantly

An `Observable.Create()` method takes a factory delegate. The factory delegate is called every time an observer subscribes to the create observable. The expected return type is the IDisposable which should be given back to the call site (see Lifetime Management).

`Observable.Generate()` exists for creating sequences. It accepts a start value and three delegates: `.Generate(startValue, shouldContinue, generateNextValue, getCurrentValue)`.

`Observable.Interval()` and `Observable.Timer()` call their observers over time.

`Observable.Start(fn)` calls `fn` repeatedly, publishing whatever it returns.

You get the idea. Check out these other ways to create observables:

    Observable.Start
    Observable.FromEventPattern
    Task.ToObservable
    Task<T>.ToObservable
    IEnumerable<T>.ToObservable
    Observable.FromAsyncPattern

## Reduce, Join, Aggregate, Map

Rx provides LINQ-esque methods over observables.

You're able to `.Where` etc and `.Subscribe` to that result. `.Select`, `.SelectMany`, etc shouldn't be a surprise to you - they turn an `IObservable<T1>` in to an `IObservable<T2>` in a far more convenient way than setting up `ISubject<T1, T2>`s which achieve the same thing.

## Summary and advice

Go and play with it. Create a `Subject`, attach some `IObserver`s, call `subject.OnNext(new MyType)` and watch your observers get called. 

Don't get obsessed with Rx.

Consider Rx when you're not in control of the input stream; Rx solves "The data is here, what are you going to do about it?", but not so much "I've finished processing this data, time to get some more".

Be suspcious when you see `async` inside observers (`Subscribe(async() => {})`). Is it doing what you think it's doing? Are you sure? Prove it.

Experiment with [back-pressure](https://www.reactivemanifesto.org/glossary#Back-Pressure) until you understand where Rx makes it hard. This is _especially_ important when you throw `async` in to the mix; what happens when something in your pipeline is `async`?
