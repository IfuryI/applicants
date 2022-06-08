export const findAscendingTag = (element, tag) => {
    while (element !== null && element !== window && element !== document.body && element.tagName !== tag) {
        element = element.parentNode;
    }
    return (element && element.tagName === tag) ? element : null;
};
