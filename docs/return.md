# returning
Returning is a way of giving the result of a function to the piece of code that called the function.

## Returning in ATOM

You can return in atom by using:

```
return (string/integer/expression/variable/array)
```

however that wont work on its own because it needs to be wrapped in a function first:

```
func getname() {
  return "erica"
}
```

and then called by something:

```
var name = getname()
print(name)
```
<sub>
Output:
erica
</sub>

### extra example:

```
func addxandy(x, y) {
  var result = [x + y]
  return result
}
```

and used by

```
var result = addxandy(5, 3)
print(result)
```

<sub>
output: 8 (in this example)
</sub>
