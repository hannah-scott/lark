# lark

A light markup language in the public domain.

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

````
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
  '''  Toggles pre
  ```  Toggles pre code
  ---  Section divider
  ***  Article divider
````

## To do

- implement code as a separate block
- handle code and pre classing
- clean up codebase
- write scripts to convert from lark to gemtext
- write scripts to convert from lark to markdown (maybe)
