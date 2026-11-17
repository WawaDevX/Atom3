# returning
<sub>ezt inkább le se fordítom</sub>

A "return" egy módja annak hogy a függvény eredményét visszaadjuk annak a kódrészletnek amely hívta a függvényt.

## "return" használata atom-ban.

így tuddod használni:

```
return (string/integer/expression/variable/array)
```

Ez azonban így nem fog működni mert először egy függvénybe kell raknunk:

```
func mianev() {
  return "erika"
}
```


Majd hívni kell valahonnan:

```
var nev = mianev()
print(nev)
```
<sub>
kimenet):
erika
</sub>

### extra példák:

```
func addxandy(x, y) {
  var eredmeny = [x + y]
  return eredmeny
}
```

És felhasználva:

```
var eredmeny = addxandy(5, 3)
print(eredmeny)
```

<sub>
Kimenet: 8 (ebben a példában)
</sub>
