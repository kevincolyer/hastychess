# hastychess

Hastychess, Copyright (C) GPLv3, 2016,2017, Kevin Colyer

A not very good version of chess, because, why not?!

Initially I coded this in Perl6 but performance was so poor I re-implemented in Go. About 10,000 speed boost, which makes it fun for deeper searches.


### Implements
- negamax alpha/beta search 
- quiescent search
- piece square table for material scores and position bonus
- xboard support
- uci support in protocol
- coloured console output

### Version
- v 1.1 "Blockhead" - vastly improved!
- v 1.01 "Annihilator!" - reimplimented search and finally understood what I was doing. Got move ordering sorted. Improved UCI and Xboard support massively. PV and statistics reporting.
- V 0.99 - first version to make it to github.

### Issues
- lots

### Todo
- impliment 3 move repeat rule so to stop cycling
- Stop rooks from dancing about
- add a degree of randomness to help with beginnings???? More aggression!
- document data structures
- Fen parsing needs refactoring to eliminate dies and panics - refactor out all panics and dies
- adbstract terminal output
- add logging facility

Refactoring
- refactor to support methods on p
- refactor and design a clearer object modle and data structure.
- refactor protocol support by moving into package
- refactor for polyglot or chess server structuring
- refactor to impliment timing possible and pondering
- refactor and remove comments of perl code and remove unused code
- improve test coverage for searches and evals and protocols

Features
- improve eval beyond pst and material score for protected and supported pieces
- xboard - undo
- transition table
- killer move
- history move
- special eval for quiese
- null move pruning
- remove limits on depth search
- increase book size considerable
- SAN notation

# Licence

### Exceptions
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
