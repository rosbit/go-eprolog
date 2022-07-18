% facts
% q(QuestionNO, AnswerChoice) :-
%	QuestionNO >= 1, QuestionNO =< 10,
%	AnswerChoice @>= a, AnswerChoice @=< d.

next(1, q(1,b), 3).                 % if choice for q #1 is b, goto q #3
next(2, q(2,b), 4).                 % if choice for q #2 is b, goto q #4
next(2, q(1,a), q(2,b), 5).         % when answering q #2, the choices for q #1 & #2 are a & b，goto q #5
next(5, q(1,a), q(2,b), q(5,a), 7). % when answering q #5, with choices a, b, a for q #1, #2 and #5，goto q #7

% Rules: 
%  Q:  Question No.
%  C:  a list of Choice like [q(1,a), q(2,b)]
%  NQ: Next Question No.
nextQuestion(Q, C, NQ) :- getNext3(Q, C, NQ), !.
nextQuestion(Q, C, NQ) :- getNext2(Q, C, NQ), !.
nextQuestion(Q, C, NQ) :- getNext1(Q, C, NQ), !.
nextQuestion(Q, _, NQ) :- NQ is Q + 1.

% for prolog implementation without member/2 functor, uncomment it.
% member(M, L) :- memberchk(M, L).

getNext1(Q, C, NQ) :-
	next(Q, C1, NQ), member(C1, C), !.
getNext2(Q, C, NQ) :- 
	next(Q, C1, C2, NQ), member(C1, C), member(C2, C), !.
getNext3(Q, C, NQ) :-
	next(Q, C1, C2, C3, NQ), member(C1, C), member(C2, C), member(C3, C), !.


