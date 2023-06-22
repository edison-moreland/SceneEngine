use "../math"

// Everything needed for bvh construction

primitive ConstructBVH
    fun apply(s: Array[Shape] ref, split_axis: USize = 0): Shape? =>
        // Cycle through every axis; 0=x, 1=y, 2=z
        let compare_shapes = compare_shape_axis(split_axis)
        let next_axis = if split_axis == 2 then 0 else split_axis + 1 end

        match s.size()
        | 1 =>
            return s(0)?
        | 2 =>
            (var left, var right) = (s(0)?, s(1)?)
            if not compare_shapes(left, right) then
                (right, left) = (left, right)
            end

            return ShapeBVH(left, right, AABB.surrounding(
                ShapeBoundingBox(left),
                ShapeBoundingBox(right)
            ))
        end

        let sorted = QuickSort[Shape](consume s, compare_shapes)

        let split_at = sorted.size()/2
        let sorted_left = sorted.slice(where to=split_at)
        let sorted_right = sorted.slice(where from=split_at)

        let left = ConstructBVH(consume sorted_left, next_axis)?
        let right = ConstructBVH(consume sorted_right, next_axis)?

        ShapeBVH(left, right, AABB.surrounding(
            ShapeBoundingBox(left),
            ShapeBoundingBox(right)
        ))

    fun compare_shape_axis(axis: USize): {(Shape, Shape): Bool} =>
        // Given an axis, return a function that compares two shape's AABBs on that axis
        {(left: Shape, right: Shape): Bool =>
            let left_box = ShapeBoundingBox(left)
            let right_box = ShapeBoundingBox(right)

            left_box.minimum(axis) < right_box.minimum(axis)
        }


class ShapeBoundingBox is ShapeVisitor[AABB]
    fun apply(s: Shape): AABB =>
        s.accept[AABB](this)

    fun visit_shape_list(s: ShapeList box): AABB =>
        if s.shapes.size() == 0 then
            return AABB.zero()
        end

        var bounding_box: (AABB | None) = None

        for shape in s.shapes.values() do
            let bb = ShapeBoundingBox(shape)
            bounding_box = match bounding_box
            | let aa: AABB => AABB.surrounding(aa, bb)
            | None =>  bb
            end
        end

        match bounding_box
        | let bb: AABB => bb
        | None => AABB.zero()
        end

    fun visit_shape_bvh(s: ShapeBVH box): AABB =>
        s.bounding_box

    fun visit_sphere(s: Sphere box): AABB =>
        AABB(s.origin - Vec3.uniform(s.radius),
             s.origin + Vec3.uniform(s.radius))