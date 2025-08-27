import { UserProfile } from "@/lib/types/profile";
import { FaBirthdayCake, FaVenusMars, FaRulerVertical, FaWeight, FaBullseye, FaCalendarAlt, FaClock } from "react-icons/fa";

type Props = { profile: UserProfile };

export default function ProfileDetails({ profile }: Props) {
  return (
    <div className="p-6 bg-white rounded-2xl shadow-xl space-y-4 hover:shadow-2xl transition-shadow duration-300">
      <h3 className="text-2xl font-bold text-gray-800 border-b pb-2 mb-4">Profile Details</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="flex items-center gap-3 p-3 bg-indigo-50 rounded-lg shadow-inner hover:bg-indigo-100 transition-colors">
          <FaBirthdayCake className="text-indigo-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Age: {profile.age}</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-pink-50 rounded-lg shadow-inner hover:bg-pink-100 transition-colors">
          <FaVenusMars className="text-pink-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Gender: {profile.gender}</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-green-50 rounded-lg shadow-inner hover:bg-green-100 transition-colors">
          <FaRulerVertical className="text-green-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Height: {profile.height} cm</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-yellow-50 rounded-lg shadow-inner hover:bg-yellow-100 transition-colors">
          <FaWeight className="text-yellow-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Weight: {profile.weight} kg</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-purple-50 rounded-lg shadow-inner hover:bg-purple-100 transition-colors">
          <FaBullseye className="text-purple-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Goals: {profile.goals}</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-orange-50 rounded-lg shadow-inner hover:bg-orange-100 transition-colors">
          <FaCalendarAlt className="text-orange-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Target Weight: {profile.target_weight} kg</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-teal-50 rounded-lg shadow-inner hover:bg-teal-100 transition-colors">
          <FaCalendarAlt className="text-teal-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Workout Days: {profile.workout_days}</span>
        </div>
        <div className="flex items-center gap-3 p-3 bg-blue-50 rounded-lg shadow-inner hover:bg-blue-100 transition-colors">
          <FaClock className="text-blue-500 w-5 h-5" />
          <span className="text-gray-700 font-medium">Preferred Time: {profile.preferred_workout_time}</span>
        </div>
      </div>
    </div>
  );
}
