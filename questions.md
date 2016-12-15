Go nitpicks

* Sorting slices is rather complicated and requires large amounts of code:

```
type PlayerRating struct {
  Player string
  Rating int
}

type ByRating []PlayerRating
func (a ByRating) Len() int           { return len(a) }
func (a ByRating) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRating) Less(i, j int) bool { return a[i].Rating < a[j].Rating }

sort.Sort(sort.Reverse(ByRating(player_scores)))
```

As opposed to sorting in C:

```

struct Rating {
	char player[512]
	int rating;
};

int compare(const void* a, const void* b)
{
	Rating* a_r = (Rating*)a;
	Rating* b_r = (Rating*)b;
	return a_r->rating < b_r->rating;
}

qsort(list, item_count, sizeof(Rating), compare);
```

or in C++:

```

struct Rating {
	char player[512]
	int rating;
};

std::sort(vector.begin(), vector.end(), [](const Rating& a, const Rating& b) { return a.rating < b.rating; });
```

* Cannot perform arithmetic in templates:
```
<table>
	<tr>
		<th>Position</th>
		<th>Player</th>
		<th>Rating</th>
	{{ range $key, $value := . }}
	<tr>
		<td>{{$key}}</td>
		<td>{{ $value.Player }}</td>
		<td>{{ $value.Rating }}</td>
	</tr>
	{{ end }}
</table>
```

To add one to the key in order to start indexing at 1 instead of 0, I had to either implement something called a FuncMap, or do what I did which was adding player rank to the struct itself. One of the creators of Go language wrote this in a thread I found (https://groups.google.com/forum/#!topic/golang-nuts/M65P28FJqkg):

```
It's a slippery slope thing. Sure, +1 is easy but then you want *2 and  then <<3 and then switch and then closures and so on. Unlike with templates in other languages, I deliberately chose not to provide a full programming language, and ask the programmer to use the one that's already there. 
```

