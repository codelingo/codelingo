#include <iostream>
#include <vector>

class Object {
public:
    Object* inner;
};

class Something {
};

int Update(Object* object, std::vector<Something*> somethings);

int main() {
    Object* obj1 = new Object();
    Object* obj2 = new Object();
    obj1->inner = obj2;

    Something* some1 = new Something();
    Something* some2 = new Something();
    Something* some3 = new Something();
    Something* some4 = new Something();

    std::vector<Something*> someVec{
            some1,
            some2,
            some3,
            some4
    };

    int objHash = Update(obj1, someVec);
    if (objHash < 10) {
        return 1;
    } else {
        return 0;
    }

}


bool Combine(Object* object, Something* something) {
    return true;
}

int Hash(Object* object) {
    return 12;
}

int Update(Object* object, std::vector<Something*> somethings) {

    for (int i = 0; i < somethings.size(); i++) {
        object = object->inner;
        if (Combine(object, somethings[i])) {
            break;
        }
    }

    return Hash(object);
};
