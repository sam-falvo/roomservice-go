roomservice-go
==============

Automatic test-case design, implementation, and coverage management tool,
inspired by the Cleanroom Software Engineering design process requirements.
(Experimental)

Why?
====

While at a recent company conference, I wondered why people don't follow
test-driven development as intended.  It's such a light-weight process that,
even when followed to the letter, you barely notice the discipline you need to
exert to follow it.  Yet, people still oppose it.

I conducted an experiment: over the course of a couple of days, I wanted to
follow text-book TDD process to evolve a reusable gap-buffer implementation,
but with the intention of being mindful of even the slightest strain or
discomfort during use.  I'll use my previous work with TDD as a control.

The first issue that I observed was surprisingly not an issue caused by duress,
but rather by euphoria.  After completing the implementation, I had what looked
like a working gap buffer, but its implementation (thanks to the "You Aren't
Gonna Need It" principle) was decidedly not truly a gap buffer.  From a user
interface point of view, it was 100% drop-in compatible with a real gap buffer,
but a gap buffer it was not!

That being said, rewriting it to use a true gap buffer implementation was
drop-dead simple.  In this respect, TDD worked beautifully well and
as-promised, allowing me to change the implementation at will so long as it
doesn't break the interface.

The next issue I observed was that, despite using test-driven techniques, I
found my test coverage was significantly less than I expected.  Again, this is
a problem with euphoria -- the excitement of seeing my software evolve and work
before my eyes encouraged me to continue with feature development, instead of
writing tests that cover known edge-cases.  Countering the euphoria was my
natural propensity for procrastination.  I _knew_ that I had edge cases to
guard against, but I would tell myself that I could write their tests after I
finish just one more feature.  Eventually, I finished the gap buffer
implementation without covering any edge cases, and by then, I'd either
forgotten what the edge cases were, or just didn't feel the motivation to
complete them.

This insight motivated my work on this project -- maybe the reason others don't
seem to follow TDD as intended is because of these unseen human factors.  Maybe
some degree of test coverage automation would be a useful adjunct to the TDD
process.  But, how can I ensure I get maximum test coverage with minimum human
planning effort?  It turns out that a different development process, called
Cleanroom Software Engineering, utilizes a method of enumerating test cases
which results in a substantially higher coverage.  Having experience with
Cleanroom, I can personally attest to its promise of delivering high-quality
software; however, I can also attest that I definitely desired automation
during the test-planning steps.

RoomService provides the automation I wish I had back then, and which I feel can help also in a TDD context as well.

More docs to follow.

