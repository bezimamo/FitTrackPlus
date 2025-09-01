import { NextResponse } from "next/server";

export async function GET() {
  const plans = [
    { id: "basic", name: "Basic Plan", price: "$10", features: ["Feature A", "Feature B"] },
    { id: "pro", name: "Pro Plan", price: "$20", features: ["Feature A", "Feature B", "Feature C"] },
    { id: "enterprise", name: "Enterprise Plan", price: "$50", features: ["Feature A", "Feature B", "Feature C", "Feature D"] },
  ];

  return NextResponse.json(plans);
}
