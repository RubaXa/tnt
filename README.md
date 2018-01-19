tnt
---
Tarantool Go connector like ORM.


### Usage

```go
// model/user.go
type User struct {
	tnt.Entry
	Id    int
	Name  string
	Email string
}

func (u *User) Codec(c *tnt.Codec) {
	c.Int(&u.id);
	c.String(&u.Name);
	c.String(&u.Email);
}

// db/box.go
type Cursor tnt.Cursor;

type Box struct {
	tnt.Box
}

func (b *Box) UsersSpace() (*UsersSpace, err) {
	return b.GetSpace("users", func () (tnt.Index, tnt.ISpace) {
		return tnt.Index("primary"), &UsersSpace{
			EmailIdx: "email",
		}
	}).(*UsersSpace)
}

func Get() (*Box, err) {
	return tnt.GetBox("namespace", func () (tnt.IBox tnt.Config) {
		return &Box{}, tnt.Config{
			Server: "localhost:1241",
		};
	}).(*Box);
}

// db/users-space.go
type UsersSpace struct {
	tnt.Space
	EmailIdx   string
}

func (us *UserSpace) NextEntry() tnt.IEntry {
	return &User{}
}

func (us *UsersSpace) SelectOne(c tnt.Cursor) (*User, error) {
	c.Iter = us.IterEq
	entry, err := us.SelectRaw(c).First()
	return entry.(*User), err;	
}

func (us *UsersSpace) Select(c tnt.Cursor) ([]*User, error) {
	list := make([]*User, 0, 1)
	err := raw.Each(func (entry tnt.Entry) {
		list = append(list, entry.(*User))
	})
	return list, err
}

// main.go
func main() {
	usersSpace, err := db.Get().UsersSpace()

	user, err := usersSpace.SelectOne(&db.Cursor{
		Key: 123,
	})

	usersByEmail, err := usersSpace.Select(&db.Cursor{
		Index: usersSpace.EmailIdx,
		Key: "xxx@yyy.zzz",
	})
}
```



### Development

 - `assist test` — run tests