### Happy Programming Language
> The language is still being heavily developed any bugs, bad coding practices etc will be corrected in due time
> This is a simple experimental programming language that am literally intendeding for it to become big :)

#### BUGS
> The parser is not well implemented yet, the evaluator has massive code repetitions and alot more but hey it works :)

#### hello world
```
print("Hello world")
```

#### exception handling
```
def InvalidCallback = "InvalidCallback"
def InvalidArgumentType = "InvalidArgumentType"

fun takesACallback(age, callback) {
    if age < 18 {
        # a one for introspection ---> improve it though
        if type(callback) != "function" {
            raise Exception(InvalidCallback, "Expected a callback")
        }

        return callback(age)
    }

    return nil
}

# we can catch exceptions
fun demoExceptionHandling(throw) {
    if type(throw) != "boolean" {
        raise Exception(InvalidArgumentType, "throw should be a boolean")
    }

    try {
         if throw {
            print(takesACallback (12, 8))
        }

        # this is a correct execution

        print(takesACallback (12, fun (age) {
            print("we are in the callback")
            return age
        }))
    } catch(error) {
        if type(error) == InvalidCallback {
            print("caught the error: {error[0]}")
        } else {
            raise error
        }
    }
}

# throw the exception
demoExceptionHandling(true)

# dont throw the exception
demoExceptionHandling(false)
```
