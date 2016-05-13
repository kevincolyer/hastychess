# hastychess

Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

A not very good version of chess, because, why not?!

Initially I coded this in Perl6 but performance was so poor I re-implemented in Go. About 10,000 speed boost, which makes it fun for deeper searches.

### Version 
- V 0.99 - first version to make it to github.

### Implements
- negamax alpha/beta search 
- quiescent search
- piece square table for material scores and position bonus
- transition table
- xboard support (mostly)

### Issues
- lots

### Todo
- xboard - complete undo and force
- icu server protocol
- board to FEN
- improve TT (currently very crude) to a zorbist hash
- search and eval
  - better initial sorting of moves
  - pv
  - null move pruning
  - killer and history
  - better eval for protected and supported pieces
  - remove limits on depth search
- lots of optimisations once quality of play improves
- algebraic notation

# Licence
# Exceptions
- file /hclibs/perftsuit.epd comes from ROCE-Roman's Own Chess Engine
- Book openings from CPW 1.1 Chess Programming Wiki example engine.

### GPLv3
    Hastychess, Chess game in Go. Supports Xboard and ICU
    Copyright (C) 2016 Kevin Colyer

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
