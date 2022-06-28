module yoshisaur/app

go 1.18

replace yoshisaur/api => ./api

require (
	github.com/lib/pq v1.10.6
	yoshisaur/api v0.0.0-00010101000000-000000000000
)
