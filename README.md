# escape-tui

Program to run an escape room for Xmas'23.

Config file is a set of "answer sets" that accept correct and partially-correct
answers via equality. Player may have a set of tags that can be used to
represent player states, items, or game states. Answer sets have a requires
(precondition), provides (postcondition), and consumes (subtractive relation
between pre- and post-) that provide predicates and transformations of player
tags.
