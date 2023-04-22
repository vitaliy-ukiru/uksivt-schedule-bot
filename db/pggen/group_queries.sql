--name: GroupByID :one
SELECT year, spec, num
FROM groups
WHERE id = pggen.arg('ID');

--name: IDByGroup :one
SELECT id::integer
FROM groups
WHERE year = pggen.arg('Year')
  AND spec = pggen.arg('Spec')
  AND num = pggen.arg('Num');

--name: SelectYears :many
SELECT DISTINCT year::integer
from public.groups
ORDER BY year;

--name: SelectSpecsForYear :many
SELECT DISTINCT spec
from public.groups
WHERE year = pggen.arg('Year')
ORDER BY spec;


--name: SelectNumsForYearAndSpec :many
SELECT num::integer
from public.groups
WHERE year = pggen.arg('Year')
  AND spec = pggen.arg('Spec')
ORDER BY spec;