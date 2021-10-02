package exceptions

// this is a listing of all system errors in the system
// find a way to inject this into the main system

// we need a way to load all the exceptions into the global scope

const (
	INVALID_OPERATION_EXCEPTION string = "InvalidOperationException"
	DIVIDE_BY_ZERO_EXCEPTION           = "DivideByZeroException"
	NAME_EXCEPTION                     = "NameException"
	NO_EXCEPTION                       = "NoException"
	ARITY_EXCEPTION                    = "ArityException"
	SYSTEM_EXCEPTION                   = "SystemException"
	INVALID_INDEX_EXCEPTION            = "InvalidIndexException"
	INVALID_NODE_EXCEPTION             = "InvalidNodeException"
	INVALID_OPERATOR_EXCEPTION         = "InvalidOperatorException"
	MODULE_IMPORT_EXCEPTION            = "ModuleImportException"
	INVALID_EXCEPTION_EXCEPTION        = "InvalidExceptionException"
	INTERNAL_RETURN_EXCEPTION          = "InternalReturnException"
	DIVISION_BY_ZERO_EXCEPTION         = "DivisionByZeroException"
	ACCESS_VIOLATION_EXCEPTION         = "AccessViolationException"
	INTERNAL_BREAK_EXCEPTION           = "InternalBreakException"
)

func LoadExceptionsToScope() []string {
	// this should all be loaded as const types
	return []string{
		INVALID_OPERATION_EXCEPTION,
		DIVIDE_BY_ZERO_EXCEPTION,
		NAME_EXCEPTION,
		NO_EXCEPTION,
		ARITY_EXCEPTION,
		SYSTEM_EXCEPTION,
		INVALID_INDEX_EXCEPTION,
		INVALID_NODE_EXCEPTION,
		INVALID_OPERATOR_EXCEPTION,
		MODULE_IMPORT_EXCEPTION,
		INVALID_EXCEPTION_EXCEPTION,
		DIVISION_BY_ZERO_EXCEPTION,
	}
}
