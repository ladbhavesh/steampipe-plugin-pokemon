# Table: pokemon_nature

Natures influence how a Pokémon's stats grow.

## Examples

### Basic info

```sql
select 
  name, 
  id 
from
 pokemon_nature 
```

### List specific nature with name as 'bold'

```sql
select 
  name, 
  id 
from
 pokemon_nature
where 
  name = 'bold';
```


### List Nature which hates sour flavor

```sql
select 
   name, 
   hates_flavor 
 from 
   pokemon_nature 
where 
  hates_flavor is not null and
  hates_flavor ->> 'name' = 'sour';
```