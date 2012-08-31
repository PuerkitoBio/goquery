// Copyright (c) 2012, Martin Angers & Contributors
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation and/or
// other materials provided with the distribution.
// * Neither the name of the author nor the names of its contributors may be used to
// endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS
// OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY
// WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Package goquery implements features similar to jQuery, including the chainable syntax,
to manipulate and query an HTML document.

It depends on Go's experimental html package, which must be installed so that it can be imported
as "exp/html". See this tutorial on how to install it accordingly: http://code.google.com/p/go-wiki/wiki/InstallingExp

It uses Cascadia as CSS selector (similar to Sizzle for jQuery). This dependency is automatically installed
when using "go get ..." to install GoQuery.

To provide chainable interface, error management is strict, and goquery panics if an invalid Cascadia selector
is used (basically the same behavior as jQuery/Sizzle/document.querySelectorAll, an error is thrown). This is
necessary since multiple return values cannot be used to allow a chainable interface.

It is hosted here, along with additional documentation in the README.md file:
https://github.com/puerkitobio/goquery
*/
package goquery

// DONE array.go : Positional Manipulation: First(), Last(), Eq(), Get(), Index(), Slice()
// TESTS filter.go : Filtering: Filter(), Not(), Has(), End()
// expand.go : "Expanding": Add(), AndSelf()
// query.go : Reflect (query) node: Is(), Contains(), HasClass()
// property.go : Inspect node: Contents(), Html(), Text(), Attr(), Val(), Length(), Size()
// traversal.go : Traversal: Find(), Children(), Parents...(), Next...(), Prev...(), Closest(), Siblings()
// iteration.go : Iteration: Each(), Map()
// type.go : Selection and Document

// TODO : Benchmarks

// TODO : Check each method, if it applies to any node or only Element nodes (Cascadia's selectors already make sure of that)

// TODO : Add the following methods:
// x Add() - Misc. Traversing
// - AndSelf() - Misc. Traversing
// x Attr() - Attributes
// x Children() - Tree Traversal
// - Closest() - Tree Traversal
// x Contains() (static function?) - Utilities - needs tests
// - Contents() (similar to Children(), but includes text and comment nodes, so Children() should filter them out) - Misc. Traversing
// x Each() - Traversing
// x End() - Misc. Traversing
// x Eq() - Filtering
// x Filter() - Filtering
// x Find() : Complete with Selection object and Node object as selectors - Tree Traversal
// x First() - Filtering
// x Get() - Node (DOM) Manipulation
// x Has() - Filtering
// x HasClass() - Attributes
// - Html() ? - Attributes
// x Index() - DOM Manipulation
// - Is() - Filtering
// x Last() - Filtering
// x Length() / Size() - jQUery property
// x Map() - Filtering
// - Next() - Tree traversal
// - NextAll() - Tree traversal
// - NextUntil() - Tree traversal
// x Not() - Filtering
// - Parent() - Tree traversal
// - Parents() - Tree traversal
// - ParentsUntil() - Tree traversal
// - Prev() - Tree traversal
// - PrevAll() - Tree traversal
// - PrevUntil() - Tree traversal
// x PushStack() ? - Internals
// - Siblings() - Tree traversal
// x Slice() - Filtering
// - Text() - DOM Manipulation
// x ToArray() Is not implemented, is Selection.Nodes
// x Unique() ? Or internally only, to remove duplicates and maintain node order? - Utilities
// - Val() ? - Attributes
