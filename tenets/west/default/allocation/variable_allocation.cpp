#include <stdlib.h>

// What is obtained using new must be disposed of using delete.
// Also what is obtained using malloc must be disposed using free.

class Obj {
    int x;
    public:
    Obj();
};

Obj::Obj() {
    x=1;
}

//  The first function is correct:
void ProperlyDeleted() {
    Obj * a = new Obj();
    // Some code...
    delete a;
}

// But this code is wrong, and should be detected as a wrong use of free and a missing
// delete:
void FreedObject() {
    Obj * b = new Obj();
    // Some code 
    free(b);
}

// This is wrong too, and should be detected as a wrong use of delete.
void DeletedArray(){
    Int[] intArray = malloc(10);
    // Some code...
    delete(intArray);
}

// New objects must be deleted or passed by reference to another function that would delete them.
void D ()
{
  ObjectOnHeap objectOnHeap = new ObjectOnHeap();
//   use object...
}

// This creates a problem. After calling delete, it is very possible that the memory has been
// reallocated to another variable somewhere else. Using it again will very likely cause a crash 
void E ()
{
  ObjectOnHeap objectOnHeap = new ObjectOnHeap();
  delete objectOnHeap
//   some code
 cout << objectOnHeap
}

// This is more complex for a tool to detect: The function returns a reference. Func1 uses the reference and ends up deleting the object whereas Func2 does not. Therefore Func2 creates a leak.

ObjectOnHeap Func ()
{
  ObjectOnHeap objectOnHeap = new ObjectOnHeap();
//   use object...
  delete objectOnHeap
//   some code
//   use objectOnHeap...
}
