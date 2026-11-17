# if feltétel

Az if feltételekkel csak akkor fut le egy adot kód rész ha a feltétel igaz.

```
// ha [feltétel] {
    // csinálj valamit
}
```

Az Atom a expressionökben `[]` használ.

**### példák**

```
var bemenet = input("Milyen színű a banán?")

print(bemenet)

if [bemenet == "sárga"] { // ha a bemenet sárga
    print("okos vagy!") // kiirja hogy okos vagy. 
}
```

<sub>kimenet (ha a bemenet sárga): okos vagy!</sub>

Vagy egy még egyszerűbb példa:

```
if [5 + 5 == 10] {
    print("Rendben")
} else {
    print("Lehetetlen")
}
```

## else feltétel

Az else feltételeket egy if mellett használjuk hogy valami akkor fusson le, ha a feltétel hamis.

```
var bemenet = input("Milyen színű a banán?")

print(bemenet)

if [bemenet == "sárga"] { // ha a felhasználó által megadott érték sárga
    print("Ügyes vagy!") // kiírja, hogy okos vagy.
} else { // hogyha nem sárga
    print("Rossz válasz!") // kiírja, hogy rossz válasz.
}
```

<sub>kimenet (ha a bemenet nem sárga): Rossz válasz!</sub>

# while loopok

A while loopokkal egy kód rész addig fut amíg egy feltétel igaz.

```
while [feltétel] {
    print("Ez addig fog futni amíg a feltétel igaz")
}
```

### példák

A `true` boolean használatával végtelen loopot hozhatsz létre.

```
while [true] {
    print("még megy")
}
```

Egy másik példa ahol kilép a program ha beírod hogy `quit`:

```
var bemenet = "hello"

while [bemenet ?= "quit"] {
    print("Fut")
    var bemenet = input()
}
```