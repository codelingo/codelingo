import random
import sys


def correct():
    while True:
        try:
            n = random.randint(0, 4)
            if n <= 0:
                raise ValueError("got {0}, expected > 0".format(n))
            elif n == 1:
                raise DeprecationWarning("{0} will be removed in next update".format(n), n)
            elif n == 2:
                raise NotImplementedError("no handler for {0}".format(n))
            else:
                raise MemoryError("not enough memory to handle > 2")
        except ValueError as error:
            print(error.args[0])
            break
        except DeprecationWarning as error:
            # Ok for now
            print("Deprecation warning: {0}".format(error.args[0]))
            return error.args[1]
        except NotImplementedError:
            # That's ok, try again
            continue
        except MemoryError as error:
            # Very bad, abort
            print(repr(error))
            sys.exit(1)


def incorrect1():
    while True:
        try:
            n = random.randint(0, 4)
            if n <= 0:
                raise ValueError("got {0}, expected > 0".format(n))
            elif n == 1:
                raise DeprecationWarning("{0} will be removed in next update".format(n), n)
            elif n == 2:
                raise NotImplementedError("no handler for {0}".format(n))
            else:
                raise MemoryError("not enough memory to handle > 2")
        except Exception as error:  # ISSUE
            # Something went wrong, abort
            print(error)
            sys.exit(1)


def incorrect2():
    while True:
        try:
            n = random.randint(0, 4)
            if n <= 0:
                raise ValueError("got {0}, expected > 0".format(n))
            elif n == 1:
                raise DeprecationWarning("{0} will be removed in next update".format(n), n)
            elif n == 2:
                raise NotImplementedError("no handler for {0}".format(n))
            else:
                raise MemoryError("not enough memory to handle > 2")
        except:  # ISSUE
            # Who knows what went wrong, abort
            print("error occurred")
            sys.exit(1)


def main():
    n = correct()
    if n is not None:
        print("Success from correct(), got {0}".format(n))

    n = incorrect1()
    if n is not None:
        print("Success from incorrect1(), got {0}".format(n))

    n = incorrect2()
    if n is not None:
        print("Success from incorrect2(), got {0}".format(n))


main()
