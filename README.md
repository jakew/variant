# Variant

A fun little package I wrote to break up mundane tests.

How of ten have you written a test that is super long and complicated, only to
realize that you need refactor the entire thing just so you can make another
similar test with one minor value different?

```go
func TestTheThing(t *testing.T) {
    client := helper.NewTestClient(t)
    user := newUser()

    user.Role = "user"
    user.Status = "active"

    // ... the rest of the test
}
```

With this package, you can add in "variations" and even variations of
variations.

```go
func TestTheThing(t *testing.T) {
    variation.New(func (v *variant.Variation) {
        client := helper.NewTestClient(t)
        user := newUser()

        shouldAccess := true

        v.Variation(variant.M{
            "when the role is user": func(v *variation.Variant) {
                user.Role = "user"
                shouldAccess = false
            },
            "when the role is admin": func(v *variation.Variant) {
                user.Role = "admin"
            },
        })

        v.Variation(variant.M{
            "when the status is active": func(v *variation.Variant) {
                user.Status = "active"
            },
            "when the status is suspended": func(v *variation.Variant) {
                user.Status = "suspended"
                shouldAccess = false
            },
        })
        
        // ... the rest of the test
        if !hadAccess && shouldAccess {
            t.Fatalf("unable to access when they should have been able to")
        } else if hadAccess && !shouldAccess {
            t.Fatalf("able to access when they should not have been able to")
        }
    })
}
```

The functional enclosure in the `variation.New` parameters isn't excuted 1 time.
It's excuted 4 times. Once for each _combination_ of the variants it is given.