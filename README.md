# hastychess

A not very good version of chess, because, why not?!

Initially I coded this in Perl6 but performace was so poor I re-implimented in Go. About 10,000 speed boost, which makes it fun for deeper searches.

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
