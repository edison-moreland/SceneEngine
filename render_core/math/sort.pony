use "collections"
use "debug"
use "format"

// Iterative quicksort: https://www.techiedelight.com/iterative-implementation-of-quicksort/
// Given an array and a compare function, sort in place iterativly
primitive QuickSort[T: Any val]
    fun apply(unsorted: Array[T] ref, compare: {(T, T): Bool}): Array[T]^ =>
        var subarrays = Array[(USize, USize)]() // Stack of start-end indexes to sort
        subarrays.push((0, unsorted.size()-1))

        try
            while subarrays.size() > 0 do
                (let sub_start, let sub_end) = subarrays.pop()?

                var pivot_index = sub_start

                for i in Range(sub_start, sub_end) do
                    if compare(unsorted(i)?, unsorted(sub_end)?) then
                        unsorted.swap_elements(i, pivot_index)?

                        pivot_index = pivot_index + 1
                    end
                end

                unsorted.swap_elements(pivot_index, sub_end)?

                // If pivot_index is 0, it would underflow to usize.max_value()
                if (not pivot_index == 0) and ((pivot_index - 1) > sub_start) then
                    subarrays.push((sub_start, pivot_index - 1))
                end

                if (pivot_index + 1) < sub_end then
                    subarrays.push((pivot_index + 1, sub_end))
                end
            end
        else
            Debug("OHFUCK" where stream=DebugErr)
        end

        consume unsorted
