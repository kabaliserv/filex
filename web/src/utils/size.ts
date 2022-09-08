const Size = {
    TB: 2 ** 40,
    GB: 2 ** 30,
    MB: 2 ** 20,
    KB: 2 ** 10,
};

export function SizeToString (value: number) {
    switch (true) {
        case value > Size.TB:
            return (value / Size.TB).toFixed(1) + " To";
        case value > Size.GB:
            return (value / Size.GB).toFixed(1) + " Go";
        case value > Size.MB:
            return (value / Size.MB).toFixed(1) + " Mo";
        case value > Size.KB:
            return (value / Size.KB).toFixed(1) + " Ko";
        default:
            return value + " Octet";
    }
}