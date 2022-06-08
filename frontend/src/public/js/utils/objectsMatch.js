export const objectsMatch = (obj, src) => {
    if (obj === null || src === null) {
        return null;
    }
    return Object.keys(src).every((key) => Object.prototype.hasOwnProperty.call(obj, key) && obj[key] === src[key]);
};
