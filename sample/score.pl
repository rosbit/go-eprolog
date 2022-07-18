% facts
score(4, q(_, a)).
score(3, q(_, b)).
score(2, q(_, c)).
score(1, q(_, d)).

% Rules: cacluate total scores
calcScores(Choices, Score) :- calcScoresN(Choices, Score).
calcScoresN([], 0).
calcScoresN([H | T], Score) :-
	score(S1, H),
	calcScoresN(T, S2),
	Score is S1 + S2.
