def correct():
    raise NotImplementedError("correct() isn't implemented; but is still correct")


def incorrect():
    raise Exception("incorrect() isn't implemented; but is still incorrect")  # ISSUE


def main():
    try:
        correct()
    except NotImplementedError as error:
        print("Caught: " + repr(error))

    try:
        incorrect()  # this will cause a crash
    except NotImplementedError as error:
        print("Caught: " + repr(error))


main()
