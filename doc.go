// Copyright (c) 2012, Martin Angers & Contributors
// All rights reserved.
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
// * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
// * Neither the name of the author nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Package goquery implements features similar to jQuery, including the chainable syntax,
to manipulate and query an HTML document.

It depends on Go's experimental html package, which must be installed so that it can be imported
as "exp/html". See this tutorial on how to install it accordingly: http://code.google.com/p/go-wiki/wiki/InstallingExp

It uses Cascadia as CSS selector (similar to Sizzle for jQuery). This dependency is automatically installed
when using "go get ..." to install GoQuery.
*/
package goquery

// Positional Filtering: First(), Last(), Eq()

// TODO : Benchmarks
// TODO : Add the following methods:
// x Add() - Misc. Traversing
// - AndSelf() - Misc. Traversing
// x Attr() - Attributes
// x Children() - Tree Traversal
// - Closest() - Tree Traversal
// x Contains() (static function?) - Utilities - needs tests
// - Contents() (similar to Children(), but includes text and comment nodes, so Children() should filter them out) - Misc. Traversing
// x Each() - Traversing
// - End() - Misc. Traversing
// x Eq() - Filtering
// x Filter() - Filtering
// x Find() : Complete with Selection object and Node object as selectors - Tree Traversal
// x First() - Filtering
// x Get() - Node (DOM) Manipulation
// - Has() - Filtering
// x HasClass() - Attributes
// - Html() ? - Attributes
// - Index() - DOM Manipulation
// - Is() - Filtering
// - Last() - Filtering
// - Length() / Size() - jQUery property
// - Map() - Filtering
// - Next() - Tree traversal
// - NextAll() - Tree traversal
// - NextUntil() - Tree traversal
// - Not() - Filtering
// - Parent() - Tree traversal
// - Parents() - Tree traversal
// - ParentsUntil() - Tree traversal
// - Prev() - Tree traversal
// - PrevAll() - Tree traversal
// - PrevUntil() - Tree traversal
// - PushStack() ? - Internals
// - Siblings() - Tree traversal
// - Slice() - Filtering
// - Text() - DOM Manipulation
// - ToArray() Is not implemented, is Selection.Nodes
// - Unique() ? Or internally only, to remove duplicates and maintain node order? - Utilities
// - Val() ? - Attributes
