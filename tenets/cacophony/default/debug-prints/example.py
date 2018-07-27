def correct(a, b):
    result = a + b
    return result


def incorrect(a, b):
    result = a + b
    print("result is {0}".format(result))  # ISSUE
    return result


def main():
    c = correct(4, 7)
    i = incorrect(3, 8)


main()
