import { writable, derived } from 'svelte/store';

export type CheckIn = {
  uuid: string;
  date: string;
  weight: number;
};

// Store for your data.
export const apiData = writable(Array<CheckIn>());

// Data transformation - sorting the data by date.
export const checkIns = derived(apiData, ($apiData) => {
    $apiData.sort(function(a,b){
        return new Date(b.date).valueOf() - new Date(a.date).valueOf();
      });
    return $apiData;
});