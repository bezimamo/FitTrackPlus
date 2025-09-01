export function safeUpper(value?: string | null) {
return (value ?? "").toString().toUpperCase();
}


export function getPlanImage(type?: string | null) {
switch ((type ?? '').toLowerCase()) {
case "fitness":
case "workout":
return "/assets/image/push-up.png";
case "diet":
return "/assets/image/diet.png";
case "physio":
return "/assets/image/daily-stretches.png";
default:
return "/assets/image/gym.png";
}
}