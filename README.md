# go-lifetime

Lifetime management helpers for wire.

## Background

An application my team develops is written in go and is fairly modular with dependencies injected through
constructor parameters. We build a single binary and, based on configuration, host various services and 
components. In order to manage this, we use [wire](https://github.com/google/wire), Google's compile-time 
dependency code generator. Wire looks at our constructors, analyzes the dependencies, and spits out code that constructs these
dependencies in the correct order. It also manages cleanup code, shutting down the components in reverse
order of instantiation.

### The problem

The one problem we have had with wire is that it doesn't handle conditional configuration. Let's take an 
example. I have a `PersonService` that takes a `PersonRepository` and there are multiple implementations
of the `PersonRepository`; an `InMemoryPersonRepository` and a `DBPersonRepository`. The 
`InMemoryPersonRepository` doesn't have any further dependencies, but the `DBPersonRepository` has a 
dependency on the database. 

Relevant issues:
* https://github.com/google/wire/issues/296
* https://github.com/google/wire/issues/225
* https://github.com/google/wire/issues/216

> Note that this is a somewhat contrived and simplified example, so there are other ways to handle the problem 
> in this particular case that aren't possible in a larger, more complex setup.

 The go code for this looks as follows:

```go
type Config struct {
    // If DSN is empty, then an in-memory repository will be used.
    DSN string	
}

func NewPersonService(repo PersonRepository) *PersonService {
    return &PersonService{
        Repo: repo,	
    }
}

type PersonService struct {
    Repo PersonRepository
}

type PersonRepository interface {
    GetPerson(id int) (*Person, error)
}

func NewInMemoryPersonRepository() *InMemoryPersonRepository {
    return &InMemoryPersonRepository{
        m: make(map[int]*Person),	
    }
}

type InMemoryPersonRepository struct {
    m map[int]*Person
}

func (r *InMemoryPersonRepository) GetPerson(id int) (*Person, error) {
    // do stuff
}

func NewDBPersonRepository(db *sql.DB) *DBPersonRepository {
    return &DBPersonRepository{
        db: db,
    }
}

type DBPersonRepository struct {
    db *sql.DB
}

func (r *DBPersonRepository) GetPerson(id int) (*Person, error) {
    // do stuff
}
```

When using wire, we'd have some wire helpers that looks a bit like the code below. Two things to note 
about this code:
1) Because `makePersonRepository` requires a `*sql.DB` regardless of whether it needs it, then `makeDB`
    is going to get called and therefore cannot error.
2) Even if we aren't going to use the `db` parameter in `makePersonRepository`, it still has to get passed
    in.

```go
func makeDB(cfg *Config) (*sql.DB, error) {
    if len(cfg.DSN) == 0 {
        return nil, nil	
    }	
    
    return sql.Open("mysql", cfg.DSN)
}

func makePersonRepository(cfg *Config, db *sql.DB) PersonRepository {
    if len(cfg.DSN) > 0 {
        return NewInMemoryPersonRepository()	
    }
    
    return NewDBPersonRepository(db)
}

func buildPersonService(cfg *Config) *PersonService {
    panic(wire.Build(
        NewPersonService,
        makePersonRepository,
        makeDB,
    ))
}
```

And wire would generate some code that looks a bit like this:

```go
func buildPersonService(cfg *Config) *PersonService {
    db, err := makeDB(cfg)
    if err != nil {
        panic(err)
    }
    repo := makePersonRepository(cfg, db)
    svc := NewPersonService(repo)
    return svc
}
```

### The Solution

Instead of working with instances, we can work with lifetimes. Many IoC containers have lifetimes built-in,
but wire does not and makes everything an eagerly created singleton. This library makes lifetimes explicit,
allows them to be lazily-loaded, and adds support for transient lifetimes.

The interface is:
```go
type Lifetime[T any] interface {
	Instance() (T, error)
}
```

and we can use it like this, replacing our functions above. One thing to note about the below code:
1) The `makeDB` function no longer needs to check the configuration and return placeholders. `makeDB`
    will only be called if the `db` instance is actually needed.

```go
func makeDB(cfg *Config) (lifetime.Lifetime[*sql.DB], error) {
    return lifetime.NewSingleton(func() (*sql.DB, error) {
        return sql.Open("mysql", cfg.DSN)
    }) 
}

func makePersonRepository(cfg *Config, dblife lifetime.Lifetime[*sql.DB]) lifetime.Lifetime[PersonRepository] {
    return lifetime.NewSingleton(func() (PersonRepository, error) {
        if len(cfg.DSN) > 0 {
            return NewInMemoryPersonRepository()
        }
        
        db, err := dblife.Instance()
        if err != nil {
            return fmt.Errorf("creating db: %w", err)	
        }       
        
        return NewDBPersonRepository(db), nil
    })
}

func makePersonService(repolife lifetime.Lifetime[PersonRepository]) lifetime.Lifetime[PersonService] {
    return lifetime.NewSingleton(func() (PersonService, error) {
        repo, err := repolife.Instance()
        if err != nil {
            return fmt.Errorf("creating repo: %w", err)	
        }
        
        return NewPersonService(repo), nil
    })
}

func buildPersonService(cfg *Config) lifetime.Lifetime[*PersonService] {
    panic(wire.Build(
        makePersonService,
        makePersonRepository
        makeDB,
    ))
}
```

## Installation

```
go get github.com/craiggwilson/go-lifetime
```
