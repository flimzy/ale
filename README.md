Ale is a minimalistic web framework written in Go.

## Why?

Why does the world need another web framework?

It doesn't.

I'm a strong believer in the Go philosophy of not using extraneous frameworks, and building everything yourself from scratch.

On the other hand, I write multiple web applications, and I don't want to duplicate all the boilerplate code in every one. So this package is simply *my* web framework scaffolding. I make no attempt to make this useful for anyone else. I publish it on GitHub for my own personal convenience, and on the off chance that someone else might also find it useful. But I make no guarantees, or necessarily even any effort, to make this product meaningful or useful for anyone other than me, in my own situation.

## Why the name?

There's a trend toward naming Go web frameworks after cocktails. I hate cocktails. But I like beer. And "Ale" is shorter than "Belgian Tripel", which is my favorite style, and seems catchier than either "Red" or "IPA", which tie it for length.

## What does it do?

Ale is designed to make writing 12-factor apps easier, and this informs many assumptions.

Ale simplifies a few tedious tasks I was tired of repeating in my web apps:

 * Reading configuration from the environment
 * Logging to stdout
 * Debugging to stdout when the 'debug' build flag is used
 * Port binding
 * Graceful shutdown

TODO:

 * Simple handling of html/template

## What does it not do?

Anything not on the list above, Ale does not attempt to do.

## What license?

Ale is released under the MIT license. See the LICENSE.txt file for details.
