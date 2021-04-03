# Overview

Set methods over strings.

# Usage

`go get github.com/wojnosystems/go-string-set`

```go
package main

import (
	"github.com/wojnosystems/go-string-set/string_set"
	"strings"
)

func main() {
  myPeople := string_set.New()
  myPeople.Add("Bob")
  myPeople.Add("Jane")
  myPeople.Add("Gary")
  myPeople.Add("Bob")
  
  if myPeople.Includes("Bob") {
  	// Bob exists in myPeople, so this code block executes
  }
  
  janesPeople := string_set.New()
  janesPeople.Add("Gary")
  janesPeople.Add("Terry")
  
  // contains: Bob, Jane
  exclusivelyMyFriends := myPeople.Subtract(janesPeople)
  
  // contains: Bob, Jane, Gary, Terry
  allOfOurFriends := myPeople.Union(janesPeople)
  
  // contains: Gary
  friendsInCommon := janesPeople.Intersection(myPeople)
  
  // contains, in any order: Gary, Jane, Bob
  asSlice := myPeople.ToSlice()
  
  // Custom Exists
  // returns true because of Gary
  myPeople.Any(func(v string) (didMatch bool) {
    return strings.HasSuffix(v, "ry")
	})
  
  // returns true, none of my friends' names end in "x"
  myPeople.None(func(v string) (didMatch bool) {
        return strings.HasSuffix(v, "x")
    })
}
```

# Interfaces

This set contains lots of interfaces to let you slice and dice how you want users to be able to utilize the Set. Of note are these 3 Interfaces:

* Immutable: contains only the read-only operations. Users will not be able to modify the set
* Mutable: contains the write-only operations. Users will be able to modify the set, but not view its contents
* Interface: contains all methods

This should allow you to safely use this set in other places in your code as necessary.
