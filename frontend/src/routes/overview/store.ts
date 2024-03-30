import { writable, derived } from 'svelte/store';

type CheckIn = {
  id: string;
  created_at: string;
  weight: number;
};

// Store for your data.
export const apiData = writable(Array<CheckIn>());

// Data transformation.
export const checkIns = derived(apiData, ($apiData) => {
  if ($apiData.length > 0){
    return $apiData;
  }
  return [];
});