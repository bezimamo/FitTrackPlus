type Props = { completion: number; isComplete: boolean };

export default function ProfileCompletion({ completion, isComplete }: Props) {
  return (
    <div className="p-6 bg-gradient-to-r from-indigo-50 to-purple-50 rounded-2xl shadow-xl hover:shadow-2xl transition-shadow duration-300 text-center">
      <h3 className="text-xl font-bold text-gray-800 mb-4">Profile Completion</h3>

      {/* Progress Bar */}
      <div className="relative w-full h-6 bg-gray-200 rounded-full overflow-hidden shadow-inner">
        <div
          className={`h-6 rounded-full transition-all duration-1000 ${
            isComplete
              ? "bg-gradient-to-r from-green-400 to-green-600"
              : "bg-gradient-to-r from-yellow-400 to-yellow-500"
          }`}
          style={{ width: `${completion}%` }}
        />
        <span className="absolute left-1/2 top-1/2 transform -translate-x-1/2 -translate-y-1/2 font-semibold text-white drop-shadow-md">
          {completion}%
        </span>
      </div>

      {/* Status Text */}
      <p
        className={`mt-4 font-medium text-lg ${
          isComplete ? "text-green-600" : "text-yellow-600"
        }`}
      >
        {isComplete ? "✔ Profile Complete" : "⚠ Profile Incomplete"}
      </p>
    </div>
  );
}
