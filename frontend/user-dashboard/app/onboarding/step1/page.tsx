'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Step1() {
  const router = useRouter();
  const [name, setName] = useState('');
  const [age, setAge] = useState('');

  return (
    <div className="max-w-md mx-auto bg-white p-6 rounded shadow">
      <h1 className="text-xl font-bold mb-4">Step 1: Basic Profile</h1>
      <input
        type="text"
        placeholder="Full Name"
        value={name}
        onChange={e => setName(e.target.value)}
        className="w-full p-2 border mb-4 rounded"
      />
      <input
        type="number"
        placeholder="Age"
        value={age}
        onChange={e => setAge(e.target.value)}
        className="w-full p-2 border mb-4 rounded"
      />
      <button
        onClick={() => router.push('/onboarding/step2')}
        className="bg-blue-600 text-white px-4 py-2 rounded"
      >
        Next
      </button>
    </div>
  );
}
