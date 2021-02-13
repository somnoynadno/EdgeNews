export function sortArrayByKey(arr, key, asc = true) {
    return arr.sort(function (a, b) {
        // @ts-ignore
        let x = a[key];
        let y = b[key];

        if (key === "CreatedAt") {
            if (asc) {
                return new Date(a[key]) - new Date(b[key]);
            } else {
                return new Date(b[key]) - new Date(a[key]);
            }
        }

        if (parseFloat(x) && parseFloat(y)) {
            return asc ? x - y : y - x;
        }
        return asc ? ('' + x).localeCompare(y) : ('' + y).localeCompare(x);
    });
}
