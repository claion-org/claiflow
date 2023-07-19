export const StringWithDefault = (d) => (s) =>  {
    if (s == null) return d;
    if (s.length == 0) return d;
    return s;
}

export function CompareDate(a, b) {
    // console.log(`CompareDate: a=${a} b=${b}`)
    // console.log(`CompareDate: a.getTime=${a.getTime()} b.getTime=${b.getTime()}`)

    if (typeof(a) == 'string') a = new Date(a);
    if (typeof(b) == 'string') b = new Date(b);

    return a.getTime() - b.getTime();
}

export function deepCompare(arg1, arg2) {
    if (Object.prototype.toString.call(arg1) === Object.prototype.toString.call(arg2)){
        if (Object.prototype.toString.call(arg1) === '[object Object]' || Object.prototype.toString.call(arg1) === '[object Array]' ){
            if (Object.keys(arg1).length !== Object.keys(arg2).length ){
                return false;
            }
            return (Object.keys(arg1).every(function(key){
                return deepCompare(arg1[key],arg2[key]);
            }));
        }
        return (arg1===arg2);
    }
    return false;
}