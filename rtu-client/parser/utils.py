def format_tag(ufmt: str) -> str:
    return ufmt.split("}")[-1]


def print_rle(l: list[str], templ: str):
    prev = l[0]
    count = 1
    for tag in l[1:]:
        if tag == prev:
            count += 1
        else:
            print(templ.format(f=prev, s=count))
            prev = tag
            count = 1
    print(templ.format(f=prev, s=count))
