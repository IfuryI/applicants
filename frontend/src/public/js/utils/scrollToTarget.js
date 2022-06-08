export const scrollToTargetAdjusted = (element) => {
    const headerOffset = 45;
    const elementPosition = element.getBoundingClientRect().top + window.scrollY;
    const offsetPosition = elementPosition - headerOffset;
    window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth',
    });
};
