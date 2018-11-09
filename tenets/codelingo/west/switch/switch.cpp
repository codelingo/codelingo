#include <iostream>

int main() {
    int a = 20;

    switch (a) {
        case 10:
        case 20:
            std::cout << "20 called\n";
        case 40:
            std::cout << "40 called\n";
            goto endswitch;
        case 60:
            std::cout << "60 called\n";
            return 0;
        case 80:
            std::cout << "80 called\n";
        case 100:
            if (false) {
                    std::cout << "100 called true\n";
                    break;
            }
        case 120:
            if (true) {
                    std::cout << "120 called true\n";
                    break;
            } else {
                    std::cout << "120 called false\n";
                    break;
            }
        case 140:
            if (false) {
                    std::cout << "140 called false\n";
                    break;
            }
            break;
        default:
            std::cout << "default called\n";
            break;
    }
    endswitch:
    return 0;
}
