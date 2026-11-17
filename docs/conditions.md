# if statements

if statements are a way to run something only if a condition is true.

```
if [ condition ] {
    // do something
}
```

Atom uses [] for expressions so make sure to use it right.

### examples

```
var input = input("what color is a banana?")

print(input)

if [input == "yellow"] { // if the user input is yellow
  print("you're smart!") // print youre smart!
}
```
<sub> output (if yellow): you're smart! </sub>

or an even simpler example:

```
if [5 + 5 == 10] {
  print("Ok")
} else {
  print("impossible")
}

```


## else statements

else statements are used next to an if statement to run something if the condition is false.

```
var input = input("what color is a banana?")

print(input)

if [input == "yellow"] { // if the user input is yellow
  print("you're smart!") // print youre smart!
} else { // otherwise if not true
  print("wrong answer!") // print wrong answer!
}
```

<sub> output (if not yellow): wrong answer! </sub>

# while loops

while loops are used to run a block of code as long as a condition is true.

```
while [condition] {
  print("this will run as long as the condition is true")
}
```

### examples:

you can make an infinite loop using the true keyword.

```
while [true] {
  print("and on")
}
```
another example that will quit if you type quit

```
var input = "hello"

while [input ?= "quit"] {
  print("running")
  var input = input()
}
```