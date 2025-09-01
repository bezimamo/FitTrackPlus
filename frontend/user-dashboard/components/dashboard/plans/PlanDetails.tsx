"use client";

import { motion, AnimatePresence } from "framer-motion";
import type { PlanResponse } from "@/lib/types/plans";
import { getPlanImage, safeUpper } from "./utils";

interface PlanDetailsProps {
  plan: PlanResponse | null;
  open: boolean;
  onClose: () => void;
  onRequestAssign: (planId: number) => void;
}

export default function PlanDetails({
  plan,
  open,
  onClose,
  onRequestAssign,
}: PlanDetailsProps) {
  if (!plan) return null;

  return (
    <AnimatePresence>
      {open && (
        <motion.div
          key="plan-details"
          className="fixed inset-0 z-50 flex"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
        >
          {/* Backdrop */}
          <div className="absolute inset-0 bg-black/40" onClick={onClose} />

          {/* Sliding Panel */}
          <motion.aside
            key="aside"
            initial={{ x: "100%" }}
            animate={{ x: 0 }}
            exit={{ x: "100%" }}
            transition={{ type: "spring", stiffness: 300, damping: 30 }}
            className="ml-auto h-full w-full sm:w-[480px] bg-white shadow-2xl p-6 overflow-y-auto rounded-l-2xl"
          >
            {/* Header Image */}
            <div className="relative w-full h-40 rounded-xl overflow-hidden">
              <img
                src={getPlanImage(plan.plan_type)}
                alt={plan.plan_type}
                className="w-full h-full object-cover"
              />
              <button
                onClick={onClose}
                className="absolute top-3 right-3 bg-black/50 text-white rounded-full p-2 hover:bg-black/70 transition"
              >
                ✕
              </button>
            </div>

            {/* Content */}
            <div className="mt-6">
              <h2 className="text-2xl font-bold text-gray-900">
                {safeUpper(plan.plan_type)} Plan
              </h2>
              <p className="mt-2 text-gray-600">{plan.description}</p>

              {/* Plan Info */}
              <ul className="mt-4 space-y-2">
                <li className="p-3 bg-gray-50 rounded-lg shadow-sm border border-gray-200">
                  <p>
                    <strong>Goal Type:</strong> {plan.goal_type}
                  </p>
                  <p>
                    <strong>Duration:</strong> {plan.duration} days
                  </p>
                  <p>
                    <strong>Active:</strong> {plan.is_active ? "Yes" : "No"}
                  </p>
                </li>
              </ul>

              {/* Actions */}
              <div className="mt-6 flex justify-end gap-3">
                <button
                  onClick={onClose}
                  className="px-4 py-2 rounded-lg bg-gray-200 text-gray-800 hover:bg-gray-300 transition"
                >
                  Close
                </button>
                <button
                  onClick={() => onRequestAssign(plan.id)}
                  className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition"
                >
                  Assign Plan
                </button>
              </div>
            </div>
          </motion.aside>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
