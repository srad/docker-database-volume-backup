const days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

export const formatDate = (goTimeString: string) => {
    const date = new Date(goTimeString);

    // Extract parts
    const weekDay = date.getDay();
    const day = date.getDate(); // Day of the month (1-31)
    const month = date.getMonth() + 1; // Month (0-11, so add 1)
    const year = date.getFullYear(); // Year
    const hour = date.getHours(); // Hours (0-23)
    const minute = date.getMinutes(); // Minutes (0-59)

    return `${days[weekDay]}, ${day}.${month}.${year} ${hour}:${minute}`;
};
