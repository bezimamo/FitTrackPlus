import { UserProfile } from "@/lib/types/profile";
import { FaEnvelope, FaHeartbeat, FaPills, FaRunning, FaAllergies } from "react-icons/fa";

type Props = { profile: UserProfile };

export default function RoleProfile({ profile }: Props) {
  return (
    <div className="p-6 bg-white rounded-2xl shadow-xl space-y-4 hover:shadow-2xl transition-shadow duration-300">
      <h3 className="text-2xl font-bold text-gray-800 border-b pb-2 mb-4">Role & Preferences</h3>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Communication Preference */}
        <div className="flex items-center gap-3 p-3 bg-indigo-50 rounded-lg shadow-inner hover:bg-indigo-100 transition-colors">
          <FaEnvelope className="text-indigo-500 w-5 h-5" />
          <span className="font-medium text-gray-700">Communication: {profile.communication_preference}</span>
        </div>

        {/* Medical History */}
        <div className="flex items-center gap-3 p-3 bg-red-50 rounded-lg shadow-inner hover:bg-red-100 transition-colors">
          <FaHeartbeat className="text-red-500 w-5 h-5" />
          <span className="font-medium text-gray-700">Medical History: {profile.medical_history}</span>
        </div>

        {/* Medications */}
        <div className="flex items-center gap-3 p-3 bg-yellow-50 rounded-lg shadow-inner hover:bg-yellow-100 transition-colors">
          <FaPills className="text-yellow-500 w-5 h-5" />
          <span className="font-medium text-gray-700">Medications: {profile.medications}</span>
        </div>

        {/* Physio Needs */}
        <div className="flex items-center gap-3 p-3 bg-green-50 rounded-lg shadow-inner hover:bg-green-100 transition-colors">
          <FaRunning className="text-green-500 w-5 h-5" />
          <span className="font-medium text-gray-700">Physio Needs: {profile.physio_needs}</span>
        </div>

        {/* Allergies */}
        <div className="flex items-center gap-3 p-3 bg-pink-50 rounded-lg shadow-inner hover:bg-pink-100 transition-colors col-span-1 md:col-span-2">
          <FaAllergies className="text-pink-500 w-5 h-5" />
          <span className="font-medium text-gray-700">Allergies: {profile.allergies}</span>
        </div>
      </div>
    </div>
  );
}
