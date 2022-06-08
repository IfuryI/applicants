export const getFormValues = (form) => {
    const formData = new FormData(form);
    // converting FormData object to json
    const data = {};
    formData.forEach((value, key) => data[key] = value);

    return data;
};

export const getModelingFormValues = (form) => {
    const formData = new FormData(form);
    // converting FormData object to json
    const data = {};
    formData.forEach((value, key) => data[key] = parseInt(value));

    return data;
};