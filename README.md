cmd-popcen
==========

Usage: `cmd-popcen <npar>`

An input file for this command looks like this:

	id,pop,lat,lon
	...
	id,pop,lat,lon

This command reads the standard input, then computes the weighted centrality of each point
using the pop(ulation) as weights, lat(itude) and lon(gitude) as geographic coordinates.

After all computations are completed, the command prints ten most central points in the file.

The implementation of this "pleasingly parallel" problem runs `<npar>` goroutines concurrently.

Below it the output from `cmd-popcen {1|2|3|4} < OHlatlon.txt` running on a 4-core laptop:

	rec 243021 npar 1
	time 6298215 ms
	central locations (top 10) 
	{OH0277478 40.3945036 -82.8012415 3.29389215311675e-06 131957.472781676}
	{OH0277479 40.3980107 -82.7975976 1.5602647041079342e-06 131957.68319384338}
	{OH0277422 40.3928225 -82.7870428 2.6004411735132237e-07 131958.80668174915}
	{OH0277477 40.3833773 -82.8043768 3.900661760269836e-06 131958.95391568067}
	{OH0277480 40.3833113 -82.7946947 7.974686265440553e-06 131959.72747345455}
	{OH0277415 40.3968999 -82.7830751 3.2072107806663093e-06 131960.0827119297}
	{OH0277466 40.3838122 -82.8116574 3.1205294082158684e-06 131960.91933460283}
	{OH0277467 40.395187 -82.8098262 4.0740245051707175e-06 131961.16369776713}
	{OH0277421 40.3864126 -82.78625 2.513759801062783e-06 131961.5980980697}
	{OH0277448 40.4061574 -82.7971989 9.708313714449369e-06 131961.8439780478}
	nrec 243021 npar 2
	time 3142629 ms
	central locations (top 10) 
	... output as above ...
	nrec 243021 npar 3
	time 2352738 ms
	central locations (top 10) 
	... output as above ...
	nrec 243021 npar 4
	time 1983088 ms
	central locations (top 10) 
	... output as above ...

Performance summary:

	cores     time(ms)     speedup
	  1        6298215       x1.00
	  2        3142629       x2.00
	  3        2352738       x2.68
	  4        1983088       x3.18
