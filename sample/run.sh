#!/bin/sh

function run() {
	echo ${*}
	echo "--------------------------------"
	eval ${*}
	echo
}

run ./p ./music.pl listen Who Music
run ./p ./music.pl listen ergou Music
run ./p ./music.pl listen Who bach
run ./p ./music.pl listen ergou bach
run ./p ./music.pl listen ergou mj

run ./p ./next-q.pl nextQuestion 1 "'[q(1,a)]'" NQ
run ./p ./next-q.pl nextQuestion 1 "'[q(1,b)]'" NQ
run ./p ./next-q.pl nextQuestion 5 "'[q(1,a), q(2,b), q(5,a)]'" NQ
run ./p ./next-q.pl nextQuestion 5 "'[q(1,a), q(2,b), q(5,b)]'" NQ
run ./p ./next-q.pl nextQuestion Q "'[q(1,a), q(2,b), q(5,a)]'" 7

run ./p ./score.pl calcScores "'[q(1,a), q(2,a), q(3,c), q(4,d)]'" Score
run ./p ./score.pl calcScores "'[q(1,a), q(2,b)]'" Score
