# lark

lark is a lightweight markup language designed for my website, HD-DN. The language is in the public domain and I exert no rights over it.

## Structure

Lark documents are lists of articles. Each article contains a number of sections, and each section is made up of blocks.

```
+--------------------------------------+
|  lark                                |
|  +--------------------------------+  |
|  |  article                       |  |
|  |  +--------------------------+  |  |
|  |  | section                  |  |  |
|  |  | +---------------------+  |  |  |
|  |  | | block               |  |  |  |
|  |  | +---------------------+  |  |  |
|  |  +--------------------------+  |  |
|  |                                |  |
|  |  +--------------------------+  |  |
|  |  | section                  |  |  |
|  |  | +---------------------+  |  |  |
|  |  | | block               |  |  |  |
|  |  | +---------------------+  |  |  |
|  |  | +---------------------+  |  |  |
|  |  | | block               |  |  |  |
|  |  | +---------------------+  |  |  |
|  |  | +---------------------+  |  |  |
|  |  | | block               |  |  |  |
|  |  | +---------------------+  |  |  |
|  |  +--------------------------+  |  |
|  +--------------------------------+  |
+--------------------------------------+
```

## Syntax

```
=    Header
-    Subheader
_    Subsubheader
+    Date
~    Author
@    Link
!    Image
>    Blockquote
*    Unordered list
:    Ordered list
'    Preformatted text
`    Code block
---  Section divider
***  Article divider
```

## To do

- write scripts to convert from lark to gemtext
- write scripts to convert from lark to markdown (maybe)
